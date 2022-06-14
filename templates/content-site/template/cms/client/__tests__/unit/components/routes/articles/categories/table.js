import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesCategoriesTable  from '../../../../../../src/components/routes/articles/categories/table.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesCategoriesTable, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
