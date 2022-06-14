import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
import { onFocusIn } from '../directives/on-focus-in';
import { storage as appState } from '../storages/app';
import Session from '../services/session';
import Store from '../services/store';
import Title from '../services/title';
import startRouter from './routes';
import App from '../components/App.sfc';

export default function( DI ) {

    // Create services
    const session = new Session( DI );
    const store = new Store( DI );
    new Title( DI );

    // Mount root component
    new App( {
        DI,
        ID: 'app',
        container: document.querySelector( 'body' ),
        directives: {
            ...directives,
            hrefto,
            onFocusIn
        },
        filters: {}
    } );

    // Load token
    store.csrftoken().then( () => {
        const app = appState( DI );
        app.tokenLoaded = true;
        app.sync();

        // Validate session (get session data)
        session.validateSession();

        // Init router
        startRouter( DI );
    } );
}
