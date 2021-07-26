import { createDI } from 'sham-ui';
import RoutesHomePage  from '../../../../../src/components/routes/home/page.sfc';
import { storage } from '../../../../../src/storages/session';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;

    const meta = renderer( RoutesHomePage, {
        DI
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
