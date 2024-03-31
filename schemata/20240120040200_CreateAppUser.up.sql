CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE COLLATION IF NOT EXISTS "case_insensitive" (locale="und-u-ks-level2", provider="icu", deterministic=false);

CREATE TABLE "app_user" (
                            "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                            "username" varchar(150) COLLATE "case_insensitive" NOT NULL UNIQUE,
                            "email" varchar(254) COLLATE "case_insensitive" NOT NULL UNIQUE,
                            "email_verified" bool NOT NULL DEFAULT false,
                            "password" varchar(128) NOT NULL,
                            "last_login" timestamp with time zone NULL,
                            "first_name" varchar(150) NOT NULL,
                            "last_name" varchar(150) NOT NULL,
                            "is_staff" boolean NOT NULL default false,
                            "is_active" boolean NOT NULL default true,
                            "date_joined" timestamp with time zone NOT NULL default CURRENT_TIMESTAMP
);
