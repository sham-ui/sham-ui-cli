// eslint-disable-next-line max-len
import RoutesArticlesListTable  from '../../../../../../src/components/routes/articles/list/table.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesListTable, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
