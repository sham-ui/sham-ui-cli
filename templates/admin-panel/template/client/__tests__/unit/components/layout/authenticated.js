import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
import LayoutAuthenticated  from '../../../../src/components/layout/authenticated.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';
import { storage } from '../../../../src/storages/session';


it( 'renders correctly', () => {
    const DI = createDI();

    storage( DI ).name = 'Test member';

    DI.bind( 'router', {
        generate: jest.fn().mockReturnValue( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/',
            addWatcher() {}
        }
    } );

    const meta = renderer( LayoutAuthenticated, {}, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
