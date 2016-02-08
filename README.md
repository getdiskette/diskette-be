# Diskette

MongoDB + REST API + Authentication + Authorization + Mail Notifications

## Status

Under heavy development.

## Roadmap

- 0.1.0
    - REST API
        - [x] `GET    /collection?st={sessionToken}&q={query}`
        - [x] `POST   /collection?st={sessionToken}           BODY={document}`
        - [x] `PUT    /collection?st={sessionToken}&q={query} BODY={partialDocument}`
        - [x] `DELETE /collection?st={sessionToken}&q={query}`

- 0.2.0
    - Authentication API
        - Unauthenticated User:
            - [ ] `Signup(email, password, lang string) (confirmationToken string, err error)`
            - [ ] `ConfirmSignup(confirmationToken string) error`
            - [ ] `ResendConfirmationMail(email, lang string) (confirmationToken string, err error)`
            - [ ] `Signin(email, password string) (sessionToken string, err error)`
            - [ ] `ForgotPasword(email, lang string) (resetToken string, err error)`
            - [ ] `ResetPassword(resetToken, newPassword string) error`
        - Authenticated User:
            - [ ] `Signout(sessionToken) error`
            - [ ] `SignoutAllSessions(sessionToken) error`
            - [ ] `ChangePassword(sessionToken, oldPassword, newPassword string) error`
            - [ ] `ChangeEmail(sessionToken, password, newEmail string) error`
        - Admin User:
            - [ ] `GetUsers() ([]User, error)`
            - [ ] `CreateUser(email, password, lang string) error`
            - [ ] `ChangeUserPassword(userId, newPassword string) error`
            - [ ] `ChangeUserEmail(userId, newEmail string) error`
            - [ ] `RemoveUsers(userIds ...string) error`
            - [ ] `SignoutAllSessions(userIds ...string) error`
            - [ ] `SuspendUsers(userIds ...string) error`
            - [ ] `UnsuspendUsers(userIds ...string) error`
            - [ ] `RemoveUnconfirmedUsers() error`

- 0.3.0
    - Authorization configuration
        - [ ] Document level access control. Example:
        ```json
        {
            "blog-post": {
                "read": true,
                "create": "session.userId != null",
                "update": "session.userId === doc.authorId || 'admin' in session.userRoles",
                "remove": "session.userId === doc.authorId || 'admin' in session.userRoles"
            }
        }
        ```

- 0.4.0
    - Mail notifications for:
        - [ ] onSignup
        - [ ] onResetPassword

- 1.0.0
    - [ ] Javascript library for usage in the browser

- 2.0.0
    - [ ] Admin webapp
    - [ ] Form generator

## License

MIT
