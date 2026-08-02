select 'fish' as animal, 10 as age, 100.1 as weight;

select * from animals
order by id;

insert into animals(name)
values ('horse');

update animals
set name = 'german shepherd'
where id = 1;

delete from animals
where id = 4;

drop table animals;

select * FROM people
order by id;

select * from emails
order by id;

select * from phones
order by id;

insert into people(first_name, last_name)
values('William', 'Daniels');

insert into emails(people_id, email_address)
values(4, 'jack_daniels@mail.ru');

select * from emails
where people_id = 4
order by id;

select p.first_name, p.last_name, e.email_address from people p
left join emails e on e.people_id = p.id
order by p.id;

select p.first_name, p.last_name, e.email_address from people p
left join emails e on e.people_id = p.id
where p.id = 1
order by p.id;

insert into phones(people_id, phone_number)
values(1, '+996-771-987-215');

select p.first_name, p.last_name, e.email_address, p2.phone_number from people p
left join emails e on e.people_id = p.id
left join phones p2 on p2.people_id = p.id
where p.id = 1
order by p.id;