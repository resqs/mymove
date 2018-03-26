import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { updateMove, loadMove } from './ducks';
import WizardPage from 'shared/WizardPage';
import MoveType from './MoveType';

export class MoveTypeWizardPage extends Component {
  componentDidMount() {
    this.props.loadMove(this.props.match.params.moveId);
  }
  handleSubmit = () => {
    const { pendingMoveType, updateMove } = this.props;
    //todo: we should make sure this move matches the redux state
    const moveId = this.props.match.params.moveId;
    if (pendingMoveType) {
      //don't update a move unless the type is selected
      updateMove(moveId, pendingMoveType);
    }
  };
  render() {
    const { pages, pageKey, pendingMoveType, currentMove } = this.props;
    const moveType =
      pendingMoveType || (currentMove && currentMove.selected_move_type);
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={moveType !== null}
      >
        <MoveType />
      </WizardPage>
    );
  }
}
MoveTypeWizardPage.propTypes = {
  updateMove: PropTypes.func.isRequired,
  pendingMoveType: PropTypes.string,
  currentMove: PropTypes.shape({
    selected_move_type: PropTypes.string,
    id: PropTypes.string,
  }),
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateMove, loadMove }, dispatch);
}
function mapStateToProps(state) {
  return state.submittedMoves;
}
export default connect(mapStateToProps, mapDispatchToProps)(MoveTypeWizardPage);
