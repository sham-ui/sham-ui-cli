import App from '../components/App.sht';

export default function( DI ) {
    new App( {
        DI,
        ID: 'app',
        container: document.querySelector( 'body' )
    } );
}
