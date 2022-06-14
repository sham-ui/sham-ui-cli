import ErrorsPanel  from '../../../src/components/ErrorsPanel.sht';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( ErrorsPanel, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
