

https://github.com/pelias/documentation/blob/master/structured-geocoding.md


pelias | cache
region | state
postalcode | zip
address | address1
country | "US"

don't use "locality" which is the city.




neighbourhood
borough

county
region


create table maintable (id integer, val string);
insert into maintable values (42, 'Hello'), (84, 'World');

create table staging (id integer, val string);
insert into staging values (42, 'New Hello'), (1337, 'Quack');

select * from maintable;

delete from maintable where id in (select id from staging);

insert into maintable select * From staging where id not in (select id from maintable);

select * from maintable;

from time import perf_counter 
t1 = perf_counter()
# code you want to test here.
print('{:4.2f} seconds'.format(float(perf_counter()-t1)) )

NPI dataset

https://www.cms.gov/Regulations-and-Guidance/Administrative-Simplification/NationalProvIdentStand/Downloads/Data_Dissemination_File-Code_Values.pdf

https://www.cms.gov/Regulations-and-Guidance/Administrative-Simplification/NationalProvIdentStand/Downloads/Data_Dissemination_File-Readme.pdf

https://www.cms.gov/Regulations-and-Guidance/Administrative-Simplification/NationalProvIdentStand/Downloads/NPPES_FOIA_Data-Elements_062007.pdf
