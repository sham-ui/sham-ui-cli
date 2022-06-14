import { createDI } from 'sham-ui';
import PagesCrudPage  from '../../../../src/components/pages/crud-page.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    const meta = renderer( PagesCrudPage, {
        DI,
        loadItems() {
            return Promise.resolve( [] );
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
