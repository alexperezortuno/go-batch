--- Create user
create user test1
    with encrypted password '123456';

-- Grant privileges
grant all privileges on database app_db to test1;

