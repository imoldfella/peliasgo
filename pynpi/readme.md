

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

https://experimentalcraft.wordpress.com/2017/11/01/how-to-make-a-postgis-tiger-geocoder-in-less-than-5-days/

# generate postgis loader script for all 50 states

SELECT Loader_Generate_Script(ARRAY['AL' , 'AK' , 'AS' , 'AZ' , 'AR' , 'CA' , 'CO' , 'MP' , 'CT' , 'DE' , 'DC' , 'FL' , 'GA' , 'GU' , 'HI' , 'ID' , 'IL' , 'IN' , 'IA' , 'KS' , 'KY' , 'LA' , 'ME' , 'MD' , 'MA' , 'MI' , 'MN' , 'MS' , 'MO' , 'MT' , 'NE' , 'NV' , 'NH' , 'NJ' , 'NM' , 'NY' , 'NC' , 'ND' , 'OH' , 'OK' , 'OR' , 'PA' , 'PR' , 'RI' , 'SC' , 'SD' , 'TN' , 'TX' , 'VI' , 'UT' , 'VT' , 'VA' , 'WA' , 'WV' , 'WI' , 'WY'], 'sh');

ALTER DATABASE gisdb SET search_path="$user",public, tiger;

psql -c "select st_x(g.geomout), st_y(g.geomout) from geocode('5 jonathan morris circle, 19063') as g"

brew services stop postgresql@14
brew services start postgresql@14