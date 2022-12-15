import { createDI } from 'sham-ui';
import Page  from '../../../src/components/Page.sfc';
import { storage } from '../../../src/storages/session';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();

    storage( DI ).sessionValidated = true;

    const meta = renderer( Page, {}, { DI } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
