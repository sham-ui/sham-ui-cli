import { createDI } from 'sham-ui';
import ContactPage  from '../../../../../src/components/routes/contact/page.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    const meta = renderer( ContactPage, {
        DI
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
