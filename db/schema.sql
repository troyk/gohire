create table users if not exists (
    id text primary key collate nocase,
    updated_at text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
    username text not null collate nocase,
    first_name text collate nocase,
    last_name text collate nocase,
    password text
);


insert into users(id,username,first_name,last_name) values('018f094d3f9f77fcbc2b735444c5017b','lt','Linus','Torvalds');
insert into users(id,username,first_name,last_name) values('018f094d3f9f72c3ba675fe3ea53346c','rp','Rob','Pike');
insert into users(id,username,first_name,last_name) values('018f094d3f9f7ef49aea12d37681fb1b','bj','Bill','Joy');
insert into users(id,username,first_name,last_name) values('018f094d3f9f7634819c1c0fad97fe7d','dr','Dennis','Ritchie');
insert into users(id,username,first_name,last_name) values('018f094d3f9f705f95f3a96127427281','kt','Ken', 'Thompson');