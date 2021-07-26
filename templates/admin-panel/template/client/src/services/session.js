import { inject } from 'sham-ui-macro/inject.macro';
import { storage as sessionStorage } from '../storages/session';
import { storage as appStorage } from '../storages/app';

export default class Session {
    @inject router;

    /** @type Store */
    @inject store;

    constructor( DI ) {
        this.DI = DI;
        this.app = appStorage( DI );
        this.data = sessionStorage( DI );
        DI.bind( 'session', this );
    }

    login( email, password ) {
        return this.store.login( { email, password } ).then( ( { Email, Name } ) => {
            Object.assign( this.data, {
                email: Email,
                name: Name
            } );
        } ).then(
            ::this.resetSessionValidation
        ).then(
            ::this.validateSession
        );
    }

    resetSessionValidation() {
        this._validateSessionPromise = undefined;
    }

    validateSession() {
        if ( undefined === this._validateSessionPromise ) {
            this._validateSessionPromise = this._validateSession();
        }
        return this._validateSessionPromise;
    }

    _validateSession() {
        this.data.sessionValidated = false;
        return this.store.validSession().then(
            ( { Email, Name, IsSuperuser } ) => {
                Object.assign( this.data, {
                    sessionValidated: true,
                    isAuthenticated: true,
                    email: Email,
                    name: Name,
                    isSuperuser: IsSuperuser
                } ).sync();
                return true;
            },
            () => {
                Object.assign( this.data, {
                    sessionValidated: true,
                    isAuthenticated: false,
                    email: '',
                    name: '',
                    isSuperuser: false
                } ).sync();
                return false;
            }
        );
    }

    logout() {
        this.store.logout().then(
            () => {

                // Reset router
                Object.assign( this.app, {
                    routerResolved: false
                } ).sync();

                // Reset cached session
                this.resetSessionValidation();

                // Go to login page
                this.router.navigate(
                    this.router.generate( 'login' )
                );
            }
        );
    }
}
