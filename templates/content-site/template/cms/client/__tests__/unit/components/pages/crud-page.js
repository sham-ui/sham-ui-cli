import { createDI } from 'sham-ui';
import PagesCrudPage  from '../../../../src/components/pages/crud-page.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    const meta = renderer( PagesCrudPage, {
        loadItems() {
            return Promise.resolve( [] );
        }
    }, {
        DI
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
