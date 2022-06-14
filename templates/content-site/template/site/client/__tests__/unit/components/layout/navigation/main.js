import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
import LayoutNavigationMain  from '../../../../../src/components/layout/navigation/main.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'router', {
        storage: {
            params: {
                page: 2
            },
            url: '',
            addWatcher() {}
        },
        generate: () => '/'
    } );

    const meta = renderer( LayoutNavigationMain, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );


it( 'toggle navbar', () => {
    const DI = createDI();
    DI.bind( 'router', {
        storage: {
            params: {
                page: 2
            },
            url: '',
            addWatcher() {}
        },
        generate: () => '/'
    } );

    const meta = renderer( LayoutNavigationMain, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    meta.component.container.querySelector( '.navbar-collapse' ).click();
    expect( meta.toJSON() ).toMatchSnapshot();
    meta.component.container.querySelector( '.navbar-collapse' ).click();
    expect( meta.toJSON() ).toMatchSnapshot();
} );
