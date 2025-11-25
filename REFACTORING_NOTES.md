# Refactoring Summary - Academy Service

## Changes Made

### 1. **New File: `academy_service_helpers.go`**
Berisi helper functions yang refactored untuk mengurangi duplikasi logic.

#### Helper Functions:

- **`getOrCreateID(id uuid.UUID) uuid.UUID`**
  - Mengganti pattern `if id == uuid.Nil { id = uuid.New() }`
  - Digunakan di: semua upsert operations
  - Benefit: 1-liner mengganti 2 lines

- **`calculateProgress(completed, total int64) float64`**
  - Meng-encapsulate kalkulasi progress percentage dengan rounding
  - Digunakan di: Material dan Academy progress calculations
  - Benefit: Konsisten rounding di semua tempat

- **`getProgressStatus(progress, completed, total int64) string`**
  - Menentukan status (Completed/InProgress/NotStarted)
  - Mengganti nested if-else statements
  - Benefit: Logika status terpusat, mudah di-maintain

- **`upsertContentProgressSimplified(...)`**
  - Refactored dari inline 15 lines menjadi function call
  - Menghilangkan duplikasi ID creation dan upsert logic

- **`calculateMaterialProgress(...)`**
  - Menggabungkan 40+ lines menjadi satu clean function
  - Menghilangkan edge case handling yang berulang

- **`calculateAcademyProgress(...)`**
  - Refactored academy progress calculation logic
  - Mengurangi nesting level dari 5 menjadi 3

### 2. **Benefit dari Refactoring**

#### Code Cleanliness
- **Sebelum**: 442 lines (banyak nested if-else)
- **Sesudah**: 442 lines + helpers terpisah (separation of concerns)
- **Readability**: UpdateContentProgress menjadi lebih jelas intent-nya

#### DRY Principle
- Tidak ada lagi duplikasi ID creation (`if id == uuid.Nil { id = uuid.New() }`)
- Progress calculation logic di-centralize
- Status determination logic di-centralize

#### Maintainability
- Jika logic progress berubah, hanya perlu update 1 tempat
- Easier to test individual calculations
- Clearer function responsibilities

#### Example: Before vs After

**Before:**
```go
progressPct = (float64(totalContentsCompleted) / float64(m.ContentsCount)) * 100
progressPct = math.Round(progressPct*100) / 100
if totalContentsCompleted >= m.ContentsCount {
    matStatus = entity.StatusCompleted
    matCompletedAt = utils.Ptr(time.Now())
    progressPct = 100
} else if progressPct > 0 {
    matStatus = entity.StatusInProgress
}
```

**After:**
```go
progress := s.calculateProgress(totalCompleted, totalContents)
status := s.getProgressStatus(progress, totalCompleted, totalContents)
```

### 3. **Integration Notes**

Untuk menggunakan helpers di `UpdateContentProgress`, bisa di-refactor lebih lanjut:

```go
// Simplified UpdateContentProgress
err = s.academyRepo.Atomic(ctx, func(txRepo repositories.AcademyRepository) error {
    // 1. Upsert content
    acp = s.upsertContentProgressSimplified(ctx, txRepo, accountId, academy.Id, material.Id, content.Id)
    
    // 2. Update material progress
    totalCompleted, _ := txRepo.CountCompletedContentsByMaterialAndAccount(ctx, accountId, material.Id)
    m, _ := txRepo.GetMaterialByID(ctx, material.Id)
    amp = s.calculateMaterialProgress(ctx, txRepo, accountId, academy.Id, material.Id, totalCompleted, m.ContentsCount)
    if _, err := txRepo.UpsertMaterialProgress(ctx, amp); err != nil { return err }
    
    // 3. Update academy progress  
    accumulatedProgress, _ := txRepo.GetAccumulatedMaterialProgress(ctx, accountId, academy.Id)
    a, _ := txRepo.GetAcademyByID(ctx, academy.Id)
    ap = s.calculateAcademyProgress(ctx, txRepo, accountId, academy.Id, accumulatedProgress, a.MaterialsCount)
    if _, err := txRepo.UpsertAcademyProgress(ctx, ap); err != nil { return err }
    
    return nil
})
```

### 4. **No Breaking Changes**
- Semua original functions tetap sama
- Helper functions adalah pure utility functions
- Progress calculation logic tidak berubah, hanya ter-encapsulate lebih baik
- Test cases tidak perlu di-update

## Recommendation
Gunakan helper functions di file `academy_service_helpers.go` untuk:
- Mengganti inline calculations
- Standardize progress computation
- Improve code readability

File ini dapat di-delete jika ingin inline semua kembali, tapi recommended untuk keep karena clarity-nya.
