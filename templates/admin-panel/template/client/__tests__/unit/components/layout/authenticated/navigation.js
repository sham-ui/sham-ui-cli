import { createDI } from 'sham-ui';
import hrefto from 'sham-ui-router/lib/href-to';
// eslint-disable-next-line max-len
import LayoutAuthenticatedNavigation  from '../../../../../src/components/layout/authenticated/navigation.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();

    DI.bind( 'router', {
        generate: jest.fn().mockReturnValue( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/',
            addWatcher() {}
        }
    } );

    const meta = renderer( LayoutAuthenticatedNavigation, {}, {
        DI,
        directives: {
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
