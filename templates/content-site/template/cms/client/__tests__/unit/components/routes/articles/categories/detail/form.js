import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesCategoriesDetailForm  from '../../../../../../../src/components/routes/articles/categories/detail/form.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesCategoriesDetailForm, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
