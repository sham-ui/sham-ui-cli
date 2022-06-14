import { inject } from 'sham-ui-macro/inject.macro';
import axios from 'axios';

export class API {
    @inject session;

    constructor( DI, { baseURL, onUnauthorized } ) {
        this.DI = DI;
        this.axios = axios.create( {
            baseURL,
            withCredentials: true
        } );
        this.axios.interceptors.request.use(
            ::this._requestInterceptor
        );
        this.axios.interceptors.response.use(
            ::this._responseInterceptor,
            ::this._responseFailInterceptor
        );

        this._onUnauthorized = onUnauthorized;
    }

    request( config ) {
        return this.axios.request( config );
    }

    _requestInterceptor( request ) {
        if ( this.token ) {
            request.headers[ 'X-CSRF-Token' ] = this.token;
        }
        return request;
    }

    _responseInterceptor( response ) {
        const token = response.headers[ 'x-csrf-token' ];
        if ( token ) {

            // Save token from response
            this.token = token;
        }

        // Return only data for success response
        return response.data;
    }

    _responseFailInterceptor( error ) {
        const { response, config, request } = error;
        if (
            this.session.data.isAuthenticated &&
            response &&
            401 === response.status
        ) {

            // Handle 401 HTTP status code for non authorized user
            this._onUnauthorized( {
                ...config,
                url: config.url.replace( config.baseURL, '/' )
            } );
        }
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
            error.response.data :
            {};
    }
}
