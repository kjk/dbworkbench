import React from 'react';
import * as action from './action.js';
import Popover from 'react-popover';


export default class QueryEditBar extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleSaveChanges = this.handleSaveChanges.bind(this);
    this.handleToggleSQLPreview = this.handleToggleSQLPreview.bind(this);

    // 1) Is there a way to move discard changes to here without using action?
    // 2) maybe move generateQuery from output to here?

    this.state = {
      isOpen: false,
      popOverText: "",
    };
  }

  togglePopover() {
    this.setState({ isOpen: !this.state.isOpen });
  }

  handleSaveChanges() {
    console.log("handleSaveChanges ");

    // TODO: must support multiple queries for multiple rows changes
    var query = this.props.generateQuery();
    action.executeQuery(query);
  }

  handleToggleSQLPreview() {
    console.log("handleSQLPreview");
    if (this.state.isOpen){
      this.setState({isOpen: false});
    } else {
      var query = this.props.generateQuery();
      query = query.split(";").join("\n");

      this.setState({
        popOverText: query,
        isOpen: true,
      });
    }
  }

  render() {
    return (
      <div id="query_edit_bar">
        <button className="discard_changes" onClick={this.props.onHandleDiscardChanges}>Discard Changes</button>
        <div className="row_number">{this.props.numberOfRowsEdited} edited rows</div>

        <Popover
          isOpen={this.state.isOpen}
          body={this.state.popOverText}
          target={"sql_preview"}
          targetElement={"sql_preview"}
          tipSize={10}>
          <button className="sql_preview" onClick={this.handleToggleSQLPreview.bind(this)}>{!this.state.isOpen ? "Show SQL Preview" : "Hide SQL Preview"}</button>
        </Popover>

        <button className="save_changes" onClick={this.handleSaveChanges.bind(this)}>Save Changes</button>
      </div>
    );
  }
}
