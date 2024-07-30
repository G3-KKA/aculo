table etl_service (
    id
    version
    app_name
    parseable_schema
)
table logs (
    id
    etl_service_id REFERENCES etl_service.id
    log_level
    priority
    info
    log_provider
    raw_data_id REFERENCES raw_data.id DEFAULT NULL ON DELETE CASCADE 
)
table raw_data (
    id
    etl_ommited_data text 
)


ex1. {
    etl_service
    (
     id: 1,
     version: 1, 
     app_name: 'schema_plus_raw_logloader', 
     parseable_schema: '[LOG_LEVEL]:\'TEXT\'[PRIORITY]:\'INTEGER\'[INFO]:\'TEXT\'[LOG_PROVIDER]:\'TEXT\''
    )
INPUT =>[INFO][1][USER: alfa2 loggined][app2.logging_service][fix alphaN error where N <0, we cannot change it, cuz its not actual text, but name::text + id::integer, and integer is strictly SiGnEd]
INPUT_SCHEMA => [LOG_LEVEL]:\'TEXT\'[PRIORITY]:\'INTEGER\'[INFO]:\'TEXT\'[LOG_PROVIDER]:\'TEXT\'[SPEC_MESSAGE]


    RESULT{
        etl_service
        (
         id: 1,
         version: 1, 
         app_name: 'schema_plus_raw_logloader', 
         parseable_schema: '[LOG_LEVEL]:\'TEXT\'[PRIORITY]:\'INTEGER\'[INFO]:\'TEXT\'[LOG_PROVIDER]:\'TEXT\''
        )
        logs(
            1
            1
            'info'
            1
            'user: alpha2| event:loggined| status:succesful_passfree'
            'app2.logging_service'
            1
        )
        raw_data(
            1
            '[fix alphaN error where N <0, we cannot change it, cuz its not actual text, but name::text + id::integer, and integer is strictly SiGnEd]'
        )




    }






}