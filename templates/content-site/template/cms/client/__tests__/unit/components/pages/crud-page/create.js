import PagesCrudPageCreate  from '../../../../../src/components/pages/crud-page/create.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( PagesCrudPageCreate, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
