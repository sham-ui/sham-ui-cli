const METHODS = [
    'get',
    'post',
    'put',
    'delete'
];

const FIXTURES = {};
const MOCKS = {};
const INTERCEPTORS = {};

function mockMethodFactory( method ) {
    return jest.fn().mockImplementation(
        ( url ) => {
            if ( !FIXTURES[ method ].hasOwnProperty( url ) ) {
                throw new Error( `Missing fixture for ${method.toUpperCase()} url = "${url}"` );
            }
            const request = {
                url,
                method,
                headers: {}
            };
            if ( INTERCEPTORS.request ) {
                INTERCEPTORS.request( request );
            }
            const { status, data, headers } = FIXTURES[ method ][ url ];
            if ( 200 === status ) {
                const response = {
                    data,
                    headers,
                    status
                };
                return Promise.resolve(
                    INTERCEPTORS.response.success( response )
                );
            } else {
                const error = {
                    config: {
                        url,
                        method,
                        baseURL: '/'
                    },
                    response: {
                        data,
                        headers,
                        status
                    }
                };
                return INTERCEPTORS.response.fail( error );
            }
        }
    );
}

METHODS.forEach(
    method => {
        MOCKS[ method ] = mockMethodFactory( method );
        FIXTURES[ method ] = {};
    }
);


export default {
    mocks: MOCKS,
    create: jest.fn().mockImplementation( () => {
        const instance  = {
            ...MOCKS,
            request( { url, method = 'get', data, params } ) {
                const args = [ url ];
                if ( undefined !== data ) {
                    args.push( data );
                }
                if ( undefined !== params ) {
                    args.push( params );
                }
                return instance[ method ].apply( null, args );
            },
            interceptors: {
                request: {
                    use( interceptor ) {
                        INTERCEPTORS.request = interceptor;
                    }
                },
                response: {
                    use( success, fail ) {
                        INTERCEPTORS.response = {
                            success,
                            fail
                        };
                    }
                }
            }
        };
        return instance;
    } ),
    use( method, url, data, status = 200, headers = {} ) {
        if ( !METHODS.includes( method ) ) {
            throw new Error( `Unknown method: "${method}"` );
        }
        FIXTURES[ method ][ url ] = {
            data,
            status,
            headers
        };
        return this;
    },

    useDefaultMocks() {
        return this;
    }
};
