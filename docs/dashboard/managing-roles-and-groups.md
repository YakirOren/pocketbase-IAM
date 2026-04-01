# Managing Roles & Groups

The Roles and Groups pages let you organize permissions and users.

## Roles

### Creating a Role

1. Click **Create** on the Roles page
2. Enter a role name (e.g., "editor", "viewer", "admin")
3. Save the role

### Attaching Policies to a Role

1. Select a role
2. Add one or more [policies](/concepts/policies)
3. Save

### Assigning a Role to Users

1. Go to the Users page
2. Select a user
3. Add the role
4. Save

## Groups

### Creating a Group

1. Click **Create** on the Groups page
2. Enter a group name (e.g., "engineering", "support")
3. Save the group

### Adding Users to a Group

1. Select a group
2. Add users
3. Save

### Attaching Policies to a Group

1. Select a group
2. Add one or more [policies](/concepts/policies)
3. Save

## Duplicate Prevention

IAM prevents duplicate assignments. You cannot assign the same role to a user twice, add the same user to a group twice, or attach the same policy to a role/group twice. The dashboard shows an error if you try.
