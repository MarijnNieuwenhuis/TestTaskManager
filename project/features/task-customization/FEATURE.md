# Feature: Task Customization with Priority Indicators

## Overview

**Status**: Draft
**Priority**: Medium
**Category**: Enhancement
**Target Release**: v2.0
**Created**: 2025-11-19
**Last Updated**: 2025-11-19

This feature enhances the task manager by adding visual priority indicators using the Eisenhower Matrix methodology. Users can assign emoticons and colors to tasks when creating them, representing urgency and importance levels, and filter tasks based on these visual attributes.

## Problem Statement

Currently, all tasks appear identical in the task list, making it difficult for users to quickly identify which tasks require immediate attention versus those that can be scheduled for later. Users need a visual system to prioritize tasks based on urgency and importance, following the proven Eisenhower Matrix methodology. Additionally, users need the ability to filter and view tasks by priority to focus on specific categories of work.

## Goals

**Primary Goal**:
- Enable users to categorize tasks using the Eisenhower Matrix framework (Urgent/Important quadrants) with visual indicators (emoticons and colors)

**Secondary Goals**:
- Provide instant task filtering by priority level and color
- Create visual distinction between task priorities for better scanning
- Automatically handle existing tasks with sensible defaults
- Maintain simplicity while adding powerful categorization

**Success Metrics**:
- Users can successfully assign priority to 100% of new tasks
- Filtering response time < 100ms (instant, client-side)
- Visual distinction is immediately clear (no user confusion)
- All existing tasks automatically receive default categorization
- Filter controls are intuitive (no documentation needed)

## Requirements

### Functional Requirements

1. **Priority Selection During Creation**: When creating a task, users must be able to select from 5 priority options (4 Eisenhower quadrants + default), each with a unique emoticon and color.

2. **Visual Task Display**: Tasks must display their assigned emoticon and be styled with their assigned color to provide immediate visual feedback in the task list.

3. **Task Filtering**: Users must be able to filter the task list by emoticon/priority and color, with filter controls displayed as buttons above the task list. Filters can be combined.

4. **Default Assignment**: New tasks created without explicit priority selection automatically receive the default priority (📋 Clipboard, Grey). Existing tasks are automatically assigned default values upon system upgrade.

5. **Immutable Priority**: Once a task is created with a priority, it cannot be changed (locked after creation). Users must delete and recreate the task to change priority.

### Non-Functional Requirements

- **Performance**: Filtering must be instant (< 100ms) using client-side filtering with JavaScript
- **Security**: No sensitive data involved; no additional security requirements beyond existing authentication
- **Scalability**: Client-side filtering scales with any reasonable number of tasks (< 1000)
- **Reliability**: Priority data persists with tasks in in-memory store; survives page reloads
- **Usability**: Priority selection must be intuitive; emoticons and colors must be immediately recognizable
- **Maintainability**: Code follows existing Go and JavaScript patterns in the codebase

## User Stories

### Story 1: Assign Priority to New Task

**As a** task manager user
**I want** to select a priority level (emoticon and color) when creating a new task
**So that** I can visually categorize my work based on urgency and importance

**Acceptance Criteria**:
- [ ] Priority selector appears in task creation form
- [ ] 5 options displayed: 🔥 Red, ⭐ Blue, ⚡ Yellow, 💡 Green, 📋 Grey
- [ ] Selected priority is shown on created task
- [ ] Default (📋 Grey) is auto-selected if no choice made
- [ ] Task displays emoticon and color styling in task list

### Story 2: Filter Tasks by Priority

**As a** task manager user
**I want** to filter tasks by priority and color
**So that** I can focus on specific categories of work (e.g., only urgent tasks)

**Acceptance Criteria**:
- [ ] Filter buttons appear above task list
- [ ] Clicking a filter button shows only matching tasks
- [ ] Multiple filters can be combined
- [ ] "Clear filters" or "Show all" option available
- [ ] Filtering is instant (< 100ms)
- [ ] Filter state is visually indicated (active/inactive)

### Story 3: View Existing Tasks with Defaults

**As a** task manager user with existing tasks
**I want** all my previous tasks to automatically receive default priority
**So that** I don't need to manually categorize old tasks

**Acceptance Criteria**:
- [ ] All existing tasks automatically assigned 📋 Grey default
- [ ] No user action required for migration
- [ ] Existing tasks appear in filtered views correctly
- [ ] No data loss occurs during migration

## Use Cases

### Use Case 1: Create Urgent Task

**Actor**: Task Manager User
**Preconditions**: User is on task list page, task creation form is visible
**Trigger**: User wants to create a new urgent and important task

**Main Flow**:
1. User types task title in input field
2. User clicks priority selector (dropdown or button group)
3. User selects 🔥 (Urgent & Important) option
4. Selected option shows red color and fire emoticon
5. User clicks "Add" button
6. Task appears in list with 🔥 emoticon and red styling
7. Task is visible in "All tasks" and when 🔥 filter is active

**Alternative Flows**:
- **Skip Priority Selection**: User doesn't select priority → System auto-assigns 📋 Grey default

**Postconditions**: Task is created with selected priority and appears in task list
**Error Handling**: If task creation fails (e.g., empty title), priority selection is preserved so user doesn't need to reselect

### Use Case 2: Filter by Important Tasks

**Actor**: Task Manager User
**Preconditions**: Task list contains tasks with various priorities
**Trigger**: User wants to see only important (but not urgent) tasks

**Main Flow**:
1. User views task list with all tasks visible
2. User clicks ⭐ (Important, Not Urgent) filter button
3. System instantly filters list to show only tasks with ⭐ priority
4. Button shows active state (highlighted/pressed appearance)
5. Task count updates to show filtered count
6. User can click button again or "Clear filters" to show all tasks

**Alternative Flows**:
- **Combine Filters**: User clicks multiple filter buttons → System shows tasks matching ANY selected filter (OR logic)
- **No Matching Tasks**: Filter results in zero tasks → System shows "No tasks match selected filters" message

**Postconditions**: Task list displays only tasks matching selected filter(s)
**Error Handling**: If filtering fails (JavaScript error), show all tasks and error message

### Use Case 3: System Upgrade Migration

**Actor**: System (automatic process)
**Preconditions**: System has existing tasks without priority data
**Trigger**: Application starts after upgrade to v2.0

**Main Flow**:
1. System checks each task in store
2. For tasks missing priority field, assigns default values:
   - Emoticon: 📋 (Clipboard)
   - Color: Grey (#6c757d)
3. Updates task records in in-memory store
4. Logs migration count to console/logs
5. System continues normal operation

**Alternative Flows**:
- **All Tasks Have Priority**: No migration needed → Skip process

**Postconditions**: All tasks have priority data; system operates normally
**Error Handling**: If migration fails for a task, log error but continue with other tasks

## Technical Context

### Dependencies

**Internal Dependencies**:
- **Task Model** (`internal/model/task.go`): Must be extended with Priority and Color fields
- **Task Store** (`internal/store/task_store.go`): No changes needed (stores Task structs as-is)
- **Task Service** (`internal/service/task_service.go`): Must handle priority validation and defaults
- **API Handler** (`internal/handler/api_handler.go`): Must accept priority in CreateTask request
- **Page Handler** (`internal/handler/page_handler.go`): Must pass priority data to templates
- **Templates** (`templates/index.html`): Must render emoticons, colors, and filter buttons
- **Stimulus Controller** (`static/js/controllers/tasks_controller.js`): Must handle priority selection and filtering

**External Dependencies**:
- Bootstrap 5.3 (already in use): For filter button styling
- Stimulus.js 3.2+ (already in use): For interactive priority selection and filtering
- No new external dependencies required

### Data Requirements

**Data Entities**:
- **Task**: Extended with two new fields:
  - `Priority` (string): Emoticon representing priority (e.g., "🔥", "⭐", "⚡", "💡", "📋")
  - `Color` (string): Hex color code (e.g., "#dc3545", "#0d6efd", "#ffc107", "#28a745", "#6c757d")

**Data Relationships**:
- Priority and Color are attributes of Task (no new relationships)

**Data Constraints**:
- Priority must be one of 5 valid emoticons: 🔥, ⭐, ⚡, 💡, 📋
- Color must be one of 7 valid hex codes: #dc3545 (red), #0d6efd (blue), #ffc107 (yellow), #28a745 (green), #6f42c1 (purple), #fd7e14 (orange), #6c757d (grey)
- Default priority is 📋 (Clipboard), default color is #6c757d (grey)
- Priority and color are immutable after task creation

**Data Volume**: No change to data volume (adds ~20 bytes per task)

**Priority-Color Mapping**:
| Eisenhower Quadrant | Emoticon | Name | Color | Hex Code |
|---------------------|----------|------|-------|----------|
| Urgent & Important | 🔥 | Fire | Red | #dc3545 |
| Important, Not Urgent | ⭐ | Star | Blue | #0d6efd |
| Urgent, Not Important | ⚡ | Lightning | Yellow | #ffc107 |
| Not Urgent, Not Important | 💡 | Light Bulb | Green | #28a745 |
| Default/Uncategorized | 📋 | Clipboard | Grey | #6c757d |

**Additional Colors Available**:
- Purple (#6f42c1) - Reserved for future use
- Orange (#fd7e14) - Reserved for future use

### APIs/Interfaces

**Modified APIs**:

**POST /api/tasks** (Create Task)
- **Request Body**:
```json
{
  "title": "Task title",
  "priority": "🔥",  // NEW: Optional, defaults to 📋
  "color": "#dc3545"  // NEW: Optional, defaults to #6c757d
}
```
- **Response** (201 Created):
```json
{
  "id": "1",
  "title": "Task title",
  "completed": false,
  "createdAt": "2025-11-19T10:00:00Z",
  "priority": "🔥",
  "color": "#dc3545"
}
```

**GET /api/tasks** (Get All Tasks)
- **Response** (200 OK):
```json
[
  {
    "id": "1",
    "title": "Task title",
    "completed": false,
    "createdAt": "2025-11-19T10:00:00Z",
    "priority": "🔥",  // NEW
    "color": "#dc3545"  // NEW
  }
]
```

**Unchanged APIs**:
- PATCH /api/tasks/{id}/toggle - No changes (priority cannot be changed)
- DELETE /api/tasks/{id} - No changes

**Internal Interfaces**:
- `TaskService.Create(title, priority, color)` - Extended to accept priority and color parameters
- Validation logic to ensure valid priority/color combinations

**Data Formats**:
- Emoticons stored as UTF-8 strings in JSON
- Colors stored as hex code strings (e.g., "#dc3545")

**Authentication/Authorization**:
- No changes - uses existing app-level authentication

### Security Considerations

**Authentication**:
- Uses existing authentication mechanism (no changes)

**Authorization**:
- All authenticated users can set priority on their tasks (no special permissions needed)

**Data Protection**:
- Priority and color are non-sensitive metadata
- No PII or sensitive information in these fields
- Standard HTTPS for data in transit

**Compliance**:
- No compliance impact (non-personal data)

**Threats & Mitigations**:
- **Invalid Priority Values**: Validated on backend before storage
- **XSS via Emoticons**: Emoticons rendered safely via Go template auto-escaping

### Scalability Considerations

- Client-side filtering scales well up to ~1000 tasks (instant performance)
- If task count grows beyond 1000, consider pagination or server-side filtering
- Priority and color fields add minimal storage overhead (~20 bytes per task)
- No database queries needed (in-memory store)

### Technical Constraints

- Priority cannot be changed after task creation (locked)
- Emoticon rendering depends on user's system fonts (Unicode support)
- Client-side filtering requires JavaScript enabled
- Limited to 5 priority options and 7 color options (simplicity constraint)

## User Experience

### User Flow

```
[Start: Task List Page]
    ↓
[User types task title]
    ↓
[User selects priority] → [Priority selector shows 5 options with emoticons/colors]
    ↓
[User clicks "Add" button]
    ↓
[Task appears with emoticon and color styling]
    ↓
[User clicks filter button (e.g., 🔥)]
    ↓
[List instantly filters to show only matching tasks]
    ↓
[User clicks "Clear filters" or another filter]
    ↓
[List updates instantly]
    ↓
[End]
```

### UI/UX Requirements

1. **Priority Selector**:
   - Displayed as radio buttons or button group
   - Shows emoticon + color name (e.g., "🔥 Red - Urgent & Important")
   - Horizontally arranged or stacked vertically
   - Default (📋 Grey) pre-selected
   - Located in task creation form, above or below title input

2. **Task Display**:
   - Emoticon displayed at start of task text
   - Task item background or border styled with assigned color
   - Color intensity: subtle (10-20% opacity) for background or solid for left border
   - Maintains readability with completed tasks (strikethrough still visible)

3. **Filter Buttons**:
   - Displayed above task list in horizontal row
   - Each button shows emoticon (e.g., "🔥", "⭐")
   - Button styled with corresponding color
   - Active state: filled/solid color
   - Inactive state: outline or ghost style
   - "Clear" or "Show All" button at end of filter row

4. **Responsive Design**:
   - Priority selector works on mobile (touch-friendly)
   - Filter buttons wrap on small screens
   - Emoticons render correctly across devices

5. **Visual Feedback**:
   - Selected priority highlights in creation form
   - Active filters show pressed/active state
   - Task count updates when filters applied
   - Smooth transitions when filtering (fade in/out)

### Accessibility

- Filter buttons have aria-label attributes (e.g., "Filter by urgent and important tasks")
- Priority selector options have descriptive labels read by screen readers
- Color is not the only indicator (emoticons provide non-color-based distinction)
- Keyboard navigation supported for priority selection and filters
- Focus indicators on interactive elements

## Success Criteria

### Definition of Done

- [x] All functional requirements implemented
- [ ] Task model extended with Priority and Color fields
- [ ] API accepts and returns priority/color in JSON
- [ ] Service layer validates priority/color values
- [ ] Templates render emoticons and color styling
- [ ] Stimulus controller handles priority selection
- [ ] Client-side filtering implemented
- [ ] Filter buttons rendered and functional
- [ ] Default priority applied to new and existing tasks
- [ ] All acceptance criteria met
- [ ] Code reviewed and approved
- [ ] Manual testing completed
- [ ] Documentation updated (README, architecture)
- [ ] Deployed and verified working

## Implementation Notes

### Phased Rollout

**Single Phase**:
- All features implemented together (not phased)
- Frontend and backend changes deployed simultaneously
- Immediate availability to all users

### Migration Requirements

**Existing Tasks**:
- On application startup (or first task fetch), check each task for priority field
- If missing: assign default priority (📋) and color (#6c757d)
- Migration happens automatically in TaskStore or TaskService initialization
- No user action required

**Rollback Plan**:
- If rollback needed, old version ignores new priority/color fields
- Tasks remain functional with just title, completed, createdAt fields
- No data loss (priority/color simply not displayed)

### Monitoring & Observability

**Metrics to Track**:
- Priority usage distribution: Count of tasks by priority type
- Filter usage: Track which filters are most commonly used
- Migration success: Log count of tasks migrated on startup

**Logging**:
- INFO: Log migration count on startup (e.g., "Migrated 15 tasks with default priority")
- DEBUG: Log priority validation failures
- ERROR: Log any unexpected priority/color values

**Alerting**:
- No specific alerts needed (non-critical feature)

## Additional Information

### Stakeholders

- **Product Owner**: TBD
- **Tech Lead**: TBD
- **Developer**: Claude (AI)
- **QA**: Manual testing by product owner

### Timeline

- **Discovery**: 2025-11-19
- **Design**: 2025-11-19 (this document)
- **Development**: TBD
- **Deployment**: TBD

### Risks & Concerns

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Emoticon rendering issues on some devices | Medium | Low | Use widely-supported Unicode emoticons; test on multiple devices |
| Client-side filtering slow with many tasks | Low | Low | Optimize filtering logic; consider server-side filtering if needed |
| Users confused by immutable priority | Low | Low | Clear messaging; allow delete/recreate workflow |
| Color accessibility issues | Medium | Low | Use emoticons as primary indicator; ensure sufficient contrast |

### Open Questions

- [x] Should priority be editable after creation? **Answer**: No, locked after creation
- [x] How many colors in palette? **Answer**: 7 colors (5 used, 2 reserved)
- [x] Client or server-side filtering? **Answer**: Client-side for instant performance

### Assumptions

- Users understand Eisenhower Matrix or can learn it quickly
- Emoticons render correctly on most modern browsers/devices
- Client-side filtering is acceptable (no pagination needed yet)
- Users prefer visual categorization over tags or labels

### Out of Scope

- Editing priority after task creation (locked behavior)
- Custom emoticons or colors (fixed palette)
- Sorting by priority (only filtering)
- Priority-based notifications or reminders
- Task delegation or sharing with priority
- Export/import with priority data
- Analytics dashboard for priority distribution
- Mobile app support (web-only for now)

### References

- **Eisenhower Matrix**: https://en.wikipedia.org/wiki/Time_management#The_Eisenhower_Method
- **Bootstrap 5.3 Colors**: https://getbootstrap.com/docs/5.3/utilities/colors/
- **Unicode Emoticons**: https://unicode.org/emoji/charts/full-emoji-list.html
- **Related Features**: [../task-manager/FEATURE.md](../task-manager/FEATURE.md)

## Change Log

| Date | Author | Changes |
|------|--------|---------|
| 2025-11-19 | Claude (Manager Skill) | Initial creation with complete requirements |

---

**Next Steps**:
- [ ] Review and refine this document with stakeholders
- [ ] Get approval to proceed
- [ ] Create architecture design (ARCHITECTURE.md via Architect skill)
- [ ] Create implementation plan (TODO.md via Tech Lead skill)
- [ ] Begin implementation (Developer skill)

---

<!-- COMPLETE -->
