# Database Migrations

This directory contains SQL migration files for the TetriON.WebServer database.

## Running Migrations

To run the migrations, execute them in order using `psql` or your PostgreSQL client:

```bash
psql -U your_username -d your_database -f migrations/001_create_users_table.sql
```

Or connect to your database and run:

```sql
\i migrations/001_create_users_table.sql
```

## Migration Files

- `001_create_users_table.sql` - Creates the users table with authentication fields
