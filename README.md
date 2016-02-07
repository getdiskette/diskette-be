# Diskette

MongoDB backend + REST frontend + Authentication + Authorization + Mail Notifications

## Status

Under heavy development.

## Roadmap

- 0.1.0
    - REST API
        - [ ] `GET /db/col?st={sessionToken}&q={query}`
        - [ ] `POST /db/col?st={sessionToken} BODY={doc}`
        - [ ] `PUT /db/col?st={sessionToken}&q={query} BODY={partialDoc}`
        - [ ] `DELETE /db/col?st={sessionToken}&q={query}`

- 0.2.0
    - Authorization Configuration
        - [ ] document level access control, example:
        ```json
        {
            "blog-post": {
                "read": true,
                "create": "session.userId != null",
                "update": "session.userId === doc._id || 'admin' in session.userRoles",
                "remove": "session.userId === doc._id || 'admin' in session.userRoles"
            }
        }
        ```

- 0.3.0
    - Authentication API
        - Unauthenticated User:
            - [ ] `Signup(email, password, lang string) (confirmationToken string, err error)`
            - [ ] `ConfirmSignup(confirmationToken string) error`
            - [ ] `ResendConfirmationMail(email, lang string) (confirmationToken string, err error)`
            - [ ] `Signin(email, password string) (sessionToken string, err error)`
            - [ ] `ForgotPasword(email, lang string) (resetToken string, err error)`
            - [ ] `ResetPassword(resetToken, newPassword string) error`
        - Authenticated User:
            - [ ] `ChangePassword(sessionToken, oldPassword, newPassword string) error`
            - [ ] `ChangeEmail(sessionToken, password, newEmail string) error`
        - Admin User:
            - [ ] `GetUsers() ([]User, error)`
            - [ ] `CreateUser(email, password, lang string) error`
            - [ ] `ChangeUserPassword(userId, newPassword string) error`
            - [ ] `ChangeUserEmail(userId, newEmail string) error`
            - [ ] `RemoveUsers(userIds ...string) error`
            - [ ] `SuspendUsers(userIds ...string) error`
            - [ ] `UnsuspendUsers(userIds ...string) error`
            - [ ] `RemoveUnconfirmedUsers() error`
    - Mail notifications for:
        - [ ] onSignup
        - [ ] onResetPassword

- 0.4.0
    - [ ] Javascript library for usage in the browser

- 1.0.0
    - [ ] Admin webapp

## License

MIT
