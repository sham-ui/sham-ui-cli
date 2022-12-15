import { createRootContext } from 'sham-ui';
import App from '../components/App.sht';

export default function( DI ) {
    new App(
        createRootContext( {
            DI,
            ID: 'app',
            container: document.querySelector( 'body' )
        } )
    );
}
