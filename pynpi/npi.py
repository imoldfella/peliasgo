
import duckdb
import pandas


# con = duckdb.connect(":memory:")
# con.execute("""
#     COPY (select * from read_csv_auto('/Users/jim/dev/asset/npi/NPPES_Data_Dissemination_September_2022/npidata_pfile_20050523-20220911.csv', header=TRUE, ALL_VARCHAR=TRUE, delim=',',NORMALIZE_NAMES=TRUE ) )
#     TO 'npi.parquet' 
#     (FORMAT 'PARQUET', CODEC 'ZSTD');
#     """)

con = duckdb.connect("db1")

# {table} = npi or update
# {table} | csv exactly
# {table}b | simplified with view
# {table}a | unique addresses concatenated with ^
# diff = updatea - npia , we need to geocode these



def export(table: str):
    con.execute(f"copy  {table} to '{table}.csv'")

def diff():
    con.execute(f"create or replace table diff as select address from updatea where address not in (select address from npia)").fetchdf()
    export('diff')

def uniqueAddress(tablename: str):
    # extract the unique addresses
    con.execute(f"create or replace table {tablename}a (address varchar primary key)") 
    con.execute(f"""insert into {tablename}a select distinct concat(substring(zip,1,5),'^',state,'^',address1)
    from {tablename}b
    where zip is not null and state is not null and address1 is not null
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
  print(con.execute(f"""select count(*) from {tablename}""").fetchall())


def whatever():
    con.execute("""select  address1,address2 from npib where address2<>'' limit 1000""").fetchdf() 
    con.execute("""select * from address order by zip,state, address1 limit 1000 """).fetchdf()
    con.execute("""select count(*) from address""").fetchall()
    con.execute("""select count(*) from npib""").fetchall()


def test1():
   load('/Users/jim/dev/asset/npi/NPPES_Data_Dissemination_092622_100222_Weekly/npidata_pfile_20220926-20221002.csv', 'update')

def test():
    load("/Users/jim/dev/asset/npi/NPPES_Data_Dissemination_September_2022/npidata_pfile_20050523-20220911.csv","npi")
    test1()
    diff()

