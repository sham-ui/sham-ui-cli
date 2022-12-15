import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesCategoriesDetail  from '../../../../../../src/components/routes/articles/categories/detail.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesCategoriesDetail, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
