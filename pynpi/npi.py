# note that paths start from root of project, not this folder.
# python3 -m pip install duckdb pandas
# shift-enter to execute a range. right-click "run selection python terminal"
import duckdb
import pandas
import os
con = duckdb.connect("build/db1")

# shared npi layer
# npi_final(tile,npi,fname,lname,sname,address_id)
# npi_address(tile,address_id,lat,lon,address1,address2,cszip_id)
# tile_cszip(tile, cszip_id, city, state, country)

# plan layer
# plan_provider(tile, contract, npi)

# plan_description (global, not per tile)
# plan_contract(code, contract, price)

# {table} = npi or update
# {table} | csv exactly
# {table}a | unique addresses concatenated with ^
# {table}b | simplified view of npi column names. used by uniqueAddress
# {table}_coded | address,lat,lon



# diff = updatea - npia , we need to geocode these

def pwd():
    print(os.getcwd())

def export(table: str):
    con.execute(f"copy  {table} to '{table}.csv'")

def diff():
    con.execute(f"create or replace table diff as select address from updatea where address not in (select address from npia)").fetchdf()
    export('build/diff')

def uniqueAddress(tablename: str):
    # extract the unique addresses
    con.execute(f"create or replace table {tablename}a (address varchar primary key)") 
    con.execute(f"""insert into {tablename}a select distinct concat(substring(zip,1,5),'^',address1)
    from {tablename}b
    where zip is not null and address1 is not null
    """)

    # create a list of addresses that we haven't seen before

# create a table and simplified view
def load(fname: str, tablename: str):
    con.execute(f"""
        create or replace table {tablename} as select * from read_csv_auto('{fname}', header=TRUE, ALL_VARCHAR=TRUE, delim=',',NORMALIZE_NAMES=TRUE ) 
        """)
    con.execute(f"""create or replace view {tablename}b as select npi,
        provider_last_name_legal_name  as lname,
        provider_first_name as fname,
        provider_credential_text as sname,
        provider_first_line_business_practice_location_address as address1,
        provider_second_line_business_practice_location_address as address2,
        provider_business_practice_location_address_city_name  as city,
        provider_business_mailing_address_postal_code as zip,
        provider_business_mailing_address_state_name as state
        from {tablename}
        """ )
    uniqueAddress(tablename)
    print(con.execute(f"select count(*) from {tablename}").fetchall())

def head(tablename: str):
  print(con.execute(f"""select  * from {tablename} limit 1000""").fetchdf())

def count(tablename: str):
  print(con.execute(f"""select count(*) from {tablename}""").fetchall())


def whatever():
    con.execute("""select  address1,address2 from npib where address2<>'' limit 1000""").fetchdf() 
    con.execute("""select * from address order by zip, address1 limit 1000 """).fetchdf()
    con.execute("""select count(*) from address""").fetchall()
    con.execute("""select count(*) from npib""").fetchall()


def test1():
   load('download/NPPES_Data_Dissemination_092622_100222_Weekly/npidata_pfile_20220926-20221002.csv', 'update')

def test():
    load("download/NPPES_Data_Dissemination_September_2022/npidata_pfile_20050523-20220911.csv","npi")
    test1()
    diff()

def loadCoded(tablename: str, csv: str):
    con.execute(f"create or replace table {tablename}(address varchar,lat double,lon double)")
    con.execute(f"copy {tablename} from '{csv}' (header false)")

loadCoded("npi_coded", "build/photon/npia2.csv")

test()
head("npi_coded")

#     con.execute(f"create or replace table coded(address varchar,lat double,lon double)")