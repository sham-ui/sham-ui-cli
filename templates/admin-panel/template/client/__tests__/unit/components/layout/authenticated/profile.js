import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
// eslint-disable-next-line max-len
import LayoutAuthenticatedProfile  from '../../../../../src/components/layout/authenticated/profile.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';
import { storage } from '../../../../../src/storages/session';

it( 'renders correctly', () => {
    const DI = createDI();

    storage( DI ).name = 'Test member';

    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/'
        }
    } );
    const meta = renderer( LayoutAuthenticatedProfile, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
