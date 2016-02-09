# Diskette

MongoDB + REST API + Authentication + Authorization + Mail Notifications

## Status

Under heavy development.

## Roadmap

- 0.1.0
    - REST API
        - [x] GET
        ```bash
        # get all documents from collection
        http localhost:5025/<mongodb_collection>
        # example:
        http localhost:5025/user
        # get the documents that match a query
        http localhost:5025/<mongodb_collection>?q='<mongodb_json_query>'
        # examples:
        http localhost:5025/user?q='{"name":"Joe Doe"}'
        http localhost:5025/user?q='{"name":{"$ne":"Joe Doe"}}'
        ```

        - [x] POST
        ```bash
        # create a new document
        http POST localhost:5025/<mongodb_collection> <mongodb_document>
        # example:
        http POST localhost:5025/user name="Joe Doe" email=joe.doe@gmail.com
        ```

        - [x] PUT
        ```bash
        # update the documents that match a query
        http PUT localhost:5025/<mongodb_collection>?q='<mongodb_json_query>' <mongodb_update>
        # example:
        http PUT localhost:5025/user?q='{"name":"Joe Doe"}' \$set:='{"email":"jdoe@gmail.com"}'
        ```

        - [x] DELETE
        ```bash
        # delete the documents that match a query
        http DELETE localhost:5025/<mongodb_collection>?q='<mongodb_json_query>'
        # example
        http DELETE localhost:5025/user?q='{"name":"Joe Doe"}'
        ```

- 0.2.0
    - Authentication API
        - Unauthenticated User:
            - [x] Signup
            ```bash
            # example
            http POST localhost:5025/user/signup \
                name="Joe Doe" email=joe.doe@gmail.com password=abc language=en
            ```
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
