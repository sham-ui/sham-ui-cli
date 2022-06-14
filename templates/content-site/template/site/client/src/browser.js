import './styles/main.scss';
import { createDI, start } from 'sham-ui';
import setupUnsafe from 'sham-ui-unsafe';
import { setup as setupRehydrator } from 'sham-ui-ssr/module/rehydrator';
import onScroll from './initializers/browser/scroll';
import enableAnimation from './initializers/browser/enable-animation';
import mainInitializer from './initializers/main';

const DI = createDI();

setupUnsafe( DI );

DI
    .bind( 'api:url', `${document.location.protocol}//${document.location.host}/api/` )
    .bind( 'location:origin', document.location.origin + '/' )
    .bind( 'location:href', document.location.href )
    .bind( 'document', document )
;

const disableRehydrating = setupRehydrator( DI, window.data );

mainInitializer(
    DI,
    document.querySelector( 'body' )
).then( () => {
    start( DI );
    disableRehydrating();
    delete window.data;

    enableAnimation();
    onScroll();
} );
