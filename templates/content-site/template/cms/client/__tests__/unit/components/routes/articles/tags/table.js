import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesTagsTable  from '../../../../../../src/components/routes/articles/tags/table.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesTagsTable, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
