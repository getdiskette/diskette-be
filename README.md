# Diskette

MongoDB backend + REST frontend + Authentication + Authorization + Mail Notifications

## Status

Under heavy development.

## Roadmap

- Authentication API
    - Unauthenticated User:
        - [ ] `Signup(email, password, lang string) (confirmationTokenStr string, err error)`
        - [ ] `ConfirmSignup(confirmationTokenStr string) error`
        - [ ] `ResendConfirmationMail(email, lang string) (confirmationTokenStr string, err error)`
        - [ ] `Signin(email, password string) (sessionTokenStr string, err error)`
        - [ ] `ForgotPasword(email, lang string) (resetTokenStr string, err error)`
        - [ ] `ResetPassword(resetTokenStr, newPassword string) error`
    - Authenticated User:
        - [ ] `ChangePassword(sessionTokenStr, oldPassword, newPassword string) error`
        - [ ] `ChangeEmail(sessionTokenStr, password, newEmail string) error`
    - Admin User:
        - [ ] `GetUsers(adminKey string) ([]User, error)`
        - [ ] `CreateUser(adminKey, email, password, lang string) error`
        - [ ] `ChangeUserPassword(adminKey, userId, newPassword string) error`
        - [ ] `ChangeUserEmail(adminKey, userId, newEmail string) error`
        - [ ] `RemoveUsers(adminKey string, userIds ...string) error`
        - [ ] `RemoveUnconfirmedUsers(adminKey string) error`
- Default mail notifications:
    - [ ] onSignup
    - [ ] onResetPassword
- REST API
    - [ ] `GET /db/col?st={sessionToken}&q={query}`
    - [ ] `POST /db/col?st={sessionToken} BODY={doc}`
    - [ ] `PUT /db/col?st={sessionToken}&q={query} BODY={partialDoc}`
    - [ ] `DELETE /db/col?st={sessionToken}&q={query}`
- Client API
    - [ ] Javascript
- Authorization Configuration
    - [ ] document level access control, example:
    ```json
    {
        "blog-post": {
            "read": true,
            "create": "userCtx.id != null",
            "update": "userCtx.id === doc._id || userCtx.role === 'admin'",
            "remove": "userCtx.id === doc._id || userCtx.role === 'admin'"
        }
    }
    ```
- [ ] Javascript client library for usage in the browser and nodejs

## License

MIT
