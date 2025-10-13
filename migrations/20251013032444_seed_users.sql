-- +goose Up
-- +goose StatementBegin
insert into users (first_name, last_name, email, birthday)
values
    ("Brandon", "Rachal", "brandon.rachal@gmail.com", "2016-07-23T00:00:00Z"),
    ("Sam", "Smith", "sam.smith@whatever.com", "2005-05-05T00:00:00Z"),
    ("John", "Doe", "john.doe@gmail.com", "2000-12-16T00:00:00Z"),
    ("Jane", "Doe", "jane.doe@gmail.com", "2003-03-16T00:00:00Z"),
    ("Ronald", "Jackson", "ronald.jackson@gmail.com", "1985-03-04T00:00:00Z"),
    ("Mary", "Perez", "mary.perez@gmail.com", "1981-05-04T00:00:00Z"),
    ("Lisa", "Mitchell", "lisa.mitchell@gmail.com", "1999-06-25T00:00:00Z"),
    ("Christopher", "Anderson", "christopher.anderson@gmail.com", "1981-02-03T00:00:00Z"),
    ("Donald", "Clark", "donald.clark@gmail.com", "1981-03-04T00:00:00Z"),
    ("Testy", "McTesterson", "testy.mctesterson@gmail.com", "1996-06-06T00:00:00Z"),
    ("Chester", "Tester", "chester.tester@gmail.com", "1965-09-07T00:00:00Z"),
    ("Mayor", "Major", "mayor.major@gmail.com", "1975-06-03T00:00:00Z"),
    ("Sally", "Supervisor", "sally.supervisor@gmail.com", "1956-02-18T00:00:00Z"),
    ("Otis", "Operator", "otis.operator@gmail.com", "1981-01-05T00:00:00Z"),
    ("Benny", "Fairfax", "benny.fairfax@gmail.com", "1943-06-09T00:00:00Z"),
    ("Allen", "Clark", "allen.clark@live.com", "1987-09-29T00:00:00Z");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from users;
-- +goose StatementEnd
