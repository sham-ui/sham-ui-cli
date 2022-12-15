import { createRootContext } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
import { storage as appState } from '../storages/app';
import Cookies from '../services/cookies';
import DarkTheme from '../services/dark-theme';
import Store from '../services/store';
import Title from '../services/title';
import SEO from '../services/seo';
import startRouter from './routes';
import formatLocaleDate from '../filters/format-locale-date';
import App from '../components/App.sfc';

export default function( DI, container ) {

    // Create services
    new Cookies( DI );
    new DarkTheme( DI, container );
    new Store( DI );
    new Title( DI );
    if ( IS_SSR ) {
        new SEO( DI );
    }

    // Mount root component
    new App(
        createRootContext( {
            DI,
            ID: 'app',
            container,
            directives: {
                ...directives,
                hrefto
            },
            filters: {
                formatLocaleDate
            }
        } )
    );

    // create storage
    appState( DI );

    // Mount router
    const routerStarted = startRouter( DI );

    return Promise.all( [
        routerStarted
    ] );
}
