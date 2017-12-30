# Google SignIn IDToken Verifier

This repository contains a go library which allows you to decode and verify Google SignIn tokens. These tokens are in fact [jwt](https://jwt.io/) tokens and fields such as expiration time, audience and issuer are checked against. JWT tokens use a signature to ensure the content has not been modified by someone else than the sender.

# Verifier type
`NewVerifier()` instantiates a new Verifier.
```go
   v := googleSignInIDTokenVerifier.NewVerifier()
```

To modify the behaviour of the verrfier, you can chain option methods.
```go
   v.LazyLoading(false).PeriodicRefresh(DefaultRefreshDuration)
```

A shared verifier instance is available at `googleSignInIDTokenVerifier.SharedInstance`
All verifier methods are then available as functions defaulting to this shared instance.
```go
   googleSignInIDTokenVerifier.Decode("...")
   // Same default instance used
   googleSignInIDTokenVerifier.SharedInstance.Decode("...")
```

# GoogleClaims
This type is returned by Decode and contains all claims in a Google SignIn jwt token.
```go
   type GoogleClaims struct {
      StandardClaims
      AuthorizedParty string `json:"azp,omitempty"`

      Name       string `json:"name,omitempty"`
      GivenName  string `json:"given_name,omitempty"`
      FamilyName string `json:"family_name,omitempty"`

      Email         string `json:"email,omitempty"`
      EmailVerified bool   `json:"email_verified,omitempty"`
      Picture       string `json:"picture,omitempty"`
      Locale        string `json:"locale,omitempty"`
   }

   type StandardClaims struct {
   	Audience  string `json:"aud,omitempty"`
   	ExpiresAt int64  `json:"exp,omitempty"`
   	Id        string `json:"jti,omitempty"`
   	IssuedAt  int64  `json:"iat,omitempty"`
   	Issuer    string `json:"iss,omitempty"`
   	NotBefore int64  `json:"nbf,omitempty"`
   	Subject   string `json:"sub,omitempty"`
   }
```

# Decode
```go
   Decode(idToken string, audience ...string) (*GoogleClaims, error)
```

If no audience is passed to Decode, no audience will be checked against.

If the token is valid, GoogleClaims are returned, in all others cases it is nil.

# Verify
```go
   Verify(idToken string, audience ...string) error
```
If no audience is passed to Decode, no audience will be checked against.

If the token is valid, error will be nil.

# Errors
These validation errors are exposed to be easier detect.
```
   ErrPublicKeyNotFound
   ErrInvalidIssuer
   ErrInvalidAudience
```

# Public Keys
Google public keys used to sign these tokens live [here](https://www.googleapis.com/oauth2/v3/certs).
Two methods exist to refresh these keys in this package.

**Lazy loading** When keys expire, the next verify or decode request will automatically refresh them. Some request can then be somewhat longer than others because an http request to Google has to be made. You should use this method if your calls to decode and verify are not latency critical or if you only call these functions a few times.

Lazy loading is the default behaviour. But you can enable it explicitly with
```go
   v.LazyLoading(true)
```

**Periodic refresh** A new go routine is created to refresh keys at periodic intervals (default 1h). You should use this method in latency critical projects or if you use this package very frequently in your project.

When enabling periodic refresh, lazy loading is automatically disabled.
```go
   v.PeriodicRefresh(time.Hour * 2)
```
To disable it once enabled, just call LazyLoading with any parameter.
```go
   v.PeriodicRefresh(time.Hour * 2)
   // Periodic refresh activated
   v.LazyLoading(true) //Either of these will work
   v.LazyLoading(false)
   // Periodic refresh deactivated
```

**Manual refresh** If you do not need any of the previous methods, you can manually refresh keys whenever you want.

To activate manual refresh, periodic refresh and lazy loading need to be deactivated.
```go
   v.LazyLoading(false)
```

You can then call `RefreshCerts` to manually scrape keys from Google. This method will refresh certs only if they are expired as set by the previous cache header.
```go
   cache, err := v.RefreshCerts()
   // cache - whether the call hit cache
```
To force keys reloading use `ForceRefreshCerts` instead.
```go
   err := v.ForceRefreshCerts()
```
