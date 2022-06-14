import axios from 'axios';

export class API {
    constructor( { baseURL } ) {
        this.axios = axios.create( {
            baseURL,
            withCredentials: true
        } );
        this.axios.interceptors.response.use(
            ::this._responseInterceptor,
            ::this._responseFailInterceptor
        );

    }

    request( config ) {
        return this.axios.request( config );
    }

    _responseInterceptor( response ) {

        // Return only data for success response
        return response.data;
    }

    _responseFailInterceptor( error ) {
        const { response, request } = error;
        if ( undefined === response && 0 === request.status ) {
            return Promise.reject( {
                Messages: [ 'Network error: connection refused. Check your network connection' ]
            } );
        }
        return Promise.reject(
            this.constructor.extractErrors( error )
        );
    }

    static extractErrors( error ) {
        return error && error.response && error.response.data ?
            { ...error.response.data, status: error.response.status } :
            {};
    }
}
