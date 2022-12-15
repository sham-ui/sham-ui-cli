import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesTagsDetail  from '../../../../../../src/components/routes/articles/tags/detail.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesTagsDetail, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
