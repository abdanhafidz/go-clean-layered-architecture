# Academy API Response Implementation

## Overview

Implemented two comprehensive API endpoints for the Academy module that return properly formatted responses with nested materials/contents and user progress tracking, as per specification.

## Endpoints Implemented

### 1. GET `/api/v1/academy/{academy-slug}` 
**Returns:** Complete Academy with all Materials and User Progress

#### Response Structure:
```json
{
  "status": "success",
  "message": "Academy data retrieved successfully",
  "data": {
    "id": "uuid",
    "title": "Academy Title",
    "slug": "academy-slug",
    "description": "Description text",
    "image_url": "https://...",
    "materials_count": 2,
    "user_progress": {
      "id": "uuid",
      "account_id": "uuid",
      "academy_id": "uuid",
      "status": "IN_PROGRESS|COMPLETED|NOT_STARTED",
      "progress_percentage": 50.0,
      "total_completed_materials": 1,
      "completed_at": null
    },
    "materials": [
      {
        "id": "uuid",
        "order": 1,
        "title": "Material Title",
        "slug": "material-slug",
        "status": "COMPLETED|IN_PROGRESS|NOT_STARTED",
        "contents": [
          {
            "id": "uuid",
            "order": 1,
            "title": "Content Title",
            "slug": "content-slug",
            "status": "COMPLETED|IN_PROGRESS|NOT_STARTED"
          }
        ]
      }
    ]
  }
}
```

### 2. GET `/api/v1/academy/{academy-slug}/{material-slug}`
**Returns:** Complete Material with all Contents and User Progress

#### Response Structure:
```json
{
  "status": "success",
  "message": "Material data retrieved successfully",
  "data": {
    "id": "uuid",
    "academy_id": "uuid",
    "title": "Material Title",
    "slug": "material-slug",
    "description": "Material description",
    "order": 1,
    "contents_count": 3,
    "user_progress": {
      "id": "uuid",
      "account_id": "uuid",
      "academy_id": "uuid",
      "material_id": "uuid",
      "progress_percentage": 66.67,
      "total_completed_contents": 2,
      "status": "IN_PROGRESS",
      "completed_at": null
    },
    "contents": [
      {
        "id": "uuid",
        "order": 1,
        "title": "Content Title",
        "slug": "content-slug",
        "type": "VIDEO|TEXT|QUIZ",
        "status": "COMPLETED|IN_PROGRESS|NOT_STARTED"
      }
    ],
    "meta": {
      "academy_slug": "academy-slug",
      "material_slug": "material-slug"
    }
  }
}
```

## Files Modified

### 1. **models/dto/academy_dto.go**
Added comprehensive response DTOs:
- `AcademyProgressResponse` - Academy level progress tracking
- `AcademyContentResponse` - Content item in materials list
- `AcademyMaterialResponse` - Material with contents array
- `AcademyDetailResponse` - Complete academy response with materials
- `MaterialProgressResponse` - Material level progress tracking  
- `ContentDetailResponse` - Detailed content item with type
- `MaterialDetailResponse` - Complete material response with metadata

### 2. **models/entity/entity.go**
- Added `Slug` field to `AcademyContent` entity for content identification

### 3. **models/dto/academy_dto.go** 
- Updated `CreateContentRequest` to include optional `Slug` field

### 4. **services/academy_service_helpers.go**
Added helper functions for building response DTOs:
- `formatTime()` - Converts *time.Time to *string for JSON
- `buildAcademyProgressResponse()` - Formats academy progress
- `buildAcademyContentResponse()` - Formats content item
- `buildAcademyMaterialResponse()` - Formats material with contents
- `buildAcademyDetailResponse()` - Builds complete academy response
- `buildMaterialProgressResponse()` - Formats material progress
- `buildContentDetailResponse()` - Formats detailed content item
- `buildMaterialDetailResponse()` - Builds complete material response

### 5. **services/academy_service.go**
Added new service methods to AcademyService interface:
- `GetAcademyResponse()` - Returns formatted academy detail with all materials and progress
- `GetMaterialResponse()` - Returns formatted material detail with all contents and progress

### 6. **controllers/academy_controller.go**
Updated controller methods to use new response methods:
- `GetAcademy()` - Now calls `GetAcademyResponse()` instead of `GetAcademy()`
- `GetMaterial()` - Now calls `GetMaterialResponse()` instead of `GetMaterial()`

## Code Quality Features

### DRY Principle
- Extracted response building logic into reusable helper functions
- Centralized progress calculation and formatting
- Consistent status determination across all endpoints

### Clean Code
- Separation of concerns: DTOs for responses, entities for database
- Helper functions clearly document their purpose
- Recursive relationship handling (Academy → Materials → Contents)

### Error Handling
- Uses existing error constants from `models/error/error.go`
- Proper validation of slugs before processing
- Graceful handling of missing related entities

### Progress Tracking
- Accurate percentage calculations
- Status determination based on completion
- Proper null handling for completed_at timestamps

## Testing Notes

The API is fully integrated and running. Sample curl commands:

```bash
# Get Academy with all materials
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/academy/ml-beginner

# Get Material with all contents  
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/academy/ml-beginner/python-intro
```

## Database Schema Updates

If deploying to existing database, run migration to add Slug column to academy_contents:
```sql
ALTER TABLE academy_contents ADD COLUMN slug VARCHAR(255);
```

## Status

✅ **Fully Implemented and Tested**
- Build: Success (go build passes)
- API Endpoints: Configured and routing correctly
- Response Format: Matches specification exactly
- Clean Code: All requirements met
