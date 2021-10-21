drop table if exists `content`;

create table `content` (
    `id` varchar(100) not null,
    `value` varchar(100) not null,
    primary key (`id`)
) engine=InnoDB default charset=utf8;

insert into content values ('1', 'v1');
insert into content values ('2', 'v2');
insert into content values ('3', 'v3');
insert into content values ('4', 'v4');

delimiter $$
create procedure getAll()
begin
    select * from content;
end $$
delimiter ;

delimiter $$
create procedure addValue(in id varchar(100), in value varchar(100))
begin
    insert into `content` values (id, value);
end $$
delimiter ;

delimiter $$
create procedure searchByValue(in ValueInTable varchar(100))
begin
    select `id`, `value` from `content` where `value` = ValueInTable ;
end $$
delimiter ;

delimiter $$
create procedure searchByIdAndValue(in IdInTable varchar(100), in ValueInTable varchar(100))
begin
    select `id`, `value` from `content` where `id` = IdInTable and `value` = ValueInTable ;
end $$
delimiter ;

delimiter $$
create procedure searchById(in IdInTable varchar(100))
begin
    select `id`, `value` from `content` where `id` = IdInTable ;
end $$
delimiter ;

delimiter $$
create procedure deleteById(in IdInTable varchar(100))
begin
    delete from `content` where `id` = IdInTable ;
end $$
delimiter ;
