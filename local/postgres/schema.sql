create table spreadsheets (
                                id serial primary key,
                                name text not null,
                                row_count int not null,
                                column_count int not null,
                                created_at timestamp default current_timestamp,
                                updated_at timestamp default current_timestamp,
                                deleted_at timestamp
);

create table cells (
                              id serial primary key,
                              spreadsheet_id int not null references spreadsheets(id),
                              raw_value text not null,
                              computed_value text,
                              row_index int not null,
                              column_index int not null,
                              created_at timestamp default current_timestamp,
                              updated_at timestamp default current_timestamp,
                              deleted_at timestamp,
                              version bigint default 1
);