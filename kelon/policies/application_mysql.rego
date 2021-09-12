package applications.mysql

# Deny all by default
allow = false

# Path: Admin can do anything
allow {
    has_role("admin")
}

# Path: GET /user/:user_id
# User can GET there own info
allow {
    some u
    input.method = "GET"
    input.path = ["user", userId]

    data.mysql.users[u].id == input.user_id
    u.id == userId
}

# Path: PUT /user/:user_id
# User can edit there own info
allow {
    some u
    input.method = "PUT"
    input.path = ["user", userId]

    data.mysql.users[u].id == input.user_id
    u.id == userId
}

# Path: GET /post/:post_id
# Everyone can read the posts as long as user is exists
allow = true {
    some user
    input.method == "GET"
    input.path = ["post", post_id]

    data.mysql.users[user].id == input.user_id
}

# Path: PUT /post/:post_id
# The owner of the post can edit it
allow = true {
    some u, p
    input.method == "PUT"
    input.path = ["post", post_id]

    # Join
    data.mysql.users[u].id == data.mysql.posts[p].owner_id

    # Where
    u.id == input.user_id
    p.id == post_id
}

# Check user has role
has_role(role_name) {
    some user
    # Query
    data.mysql.users[user].id == input.user_id
    user.role == role_name
}
