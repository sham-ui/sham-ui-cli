import { createDI } from 'sham-ui';
import hrefto from 'sham-ui-router/lib/href-to';
import LayoutNavigationLogo  from '../../../../../src/components/layout/navigation/logo.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'router', {
        generate: () => '/'
    } );

    const meta = renderer( LayoutNavigationLogo, {
        DI,
        directives: {
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
