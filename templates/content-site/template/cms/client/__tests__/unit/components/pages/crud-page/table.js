import PagesCrudPageTable  from '../../../../../src/components/pages/crud-page/table.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( PagesCrudPageTable, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
