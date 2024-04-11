import './styles/main.scss';
import { createDI, start } from 'sham-ui';
import setupUnsafe from 'sham-ui-unsafe';
import onScroll from './initializers/browser/scroll';
import enableAnimation from './initializers/browser/enable-animation';
import mainInitializer from './initializers/main';

const DI = createDI();

setupUnsafe( DI );

DI
    .bind( 'api:url', 'http://localhost:3001/api/' )
    .bind( 'location:origin', document.location.origin + '/' )
    .bind( 'location:href', document.location.href )
    .bind( 'document', document )
;

mainInitializer(
    DI,
    document.querySelector( 'body' )
).then( () => {
    start( DI );
    enableAnimation();
    onScroll();
} );

