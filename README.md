# Diskette

MongoDB + REST API + Authentication + Authorization + Mail Notifications

## Status

Under heavy development.

## Roadmap

- [x] REST

    - [x] GET
    ```bash
    # examples:
    http localhost:5025/user
    http localhost:5025/user?q='{"name":"Joe Doe"}'
    http localhost:5025/user?q='{"name":{"$ne":"Joe Doe"}}'
    ```

    - [x] POST
    ```bash
    # example:
    http POST localhost:5025/user name="Joe Doe" email=joe.doe@gmail.com
    ```

    - [x] PUT
    ```bash
    # example:
    http PUT localhost:5025/user?q='{"name":"Joe Doe"}' \$set:='{"email":"jdoe@gmail.com"}'
    ```

    - [x] DELETE
    ```bash
    # example
    http DELETE localhost:5025/user?q='{"name":"Joe Doe"}'
    ```

- [ ] Authentication

    - [x] sign up
    ```bash
    # example
    http POST localhost:5025/user/signup name="Joe Doe" email=joe.doe@gmail.com password=abc language=en
    ```

    - [x] confirm sign up
    ```bash
    http POST localhost:5025/user/confirm token=<confirmation_token>
    ```

    - [x] sign in
    ```bash
    # example
    http POST localhost:5025/user/signin email=joe.doe@gmail.com password=abc
    ```

    - [x] forgot password
    ```bash
    # example
    http POST localhost:5025/user/forgot-password email=joe.doe@gmail.com
    ```

    - [x] reset password
    ```bash
    # example
    http POST localhost:5025/user/reset-password token=<reset_token> password=123
    ```

    - [ ] `Signout(sessionToken) error`
    - [ ] `ChangePassword(sessionToken, oldPassword, newPassword string) error`
    - [ ] `UpdateProfile(sessionToken, password, newName string, newEmail string) error`

- [ ] User Management
    - [ ] `GetUsers() ([]User, error)`
    - [ ] `CreateUser(email, password, language string) error`
    - [ ] `ChangeUserPassword(userId, newPassword string) error`
    - [ ] `UpdateUserProfile(userId, newName, newEmail string) error`
    - [ ] `RemoveUsers(userIds ...string) error`
    - [ ] `Signout(userIds ...string) error`
    - [ ] `SuspendUsers(userIds ...string) error`
    - [ ] `UnsuspendUsers(userIds ...string) error`
    - [ ] `RemoveUnconfirmedUsers() error`
    - [ ] `RemoveExpiredResetKeys() error`

- [ ] Authorization
    - [ ] Document level access control
    ```json
    // example
    {
        "blog-post": {
            "read": true,
            "create": "session.userId != null",
            "update": "session.userId === doc.authorId || 'admin' in session.userRoles",
            "remove": "session.userId === doc.authorId || 'admin' in session.userRoles"
        }
    }
    ```

- [ ] Mail Notifications:
    - [ ] onSignup
    - [ ] onResetPassword

- [ ] Javascript library for usage in the browser


## License

MIT
